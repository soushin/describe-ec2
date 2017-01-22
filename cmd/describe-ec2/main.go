package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/subcommands"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&tagCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

type tagCmd struct {
	credential credential
	tagKey     string
	region     string
}

type credential struct {
	profile  string
	filename string
}

func (*tagCmd) Name() string {
	return "tag"
}

func (*tagCmd) Synopsis() string {
	return "Fetch the ec2 instance public dns name by tag search, then that stored to text file in the current directory."
}

func (*tagCmd) Usage() string {
	return `tag [-credential-profile default] [-credential-filename '~/.aws/credentials'] [-region ap-northeast-1] [-tag-key Name] '*dev*' :
  Created or updated text file
`
}

func (p *tagCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.credential.filename, "credential-filename", "", "optional: aws credential file name, when filename is empty, that will use '$HOME/.aws/credentials'")
	f.StringVar(&p.credential.profile, "credential-profile", "default", "optional: aws credential profile, default value is 'default'")
	f.StringVar(&p.region, "region", "ap-northeast-1", "optional: aws region, default value is 'ap-northeast-1'")
	f.StringVar(&p.tagKey, "tag-key", "Name", "target tag key, default value is 'Name'")
}

func (p *tagCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	text, err := getText(f.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:%s, args is %v", err.Error(), f.Args())
		fmt.Println()
		return subcommands.ExitFailure
	}

	reservations, err := getReservations(text, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:%s, search text is '%v'", err.Error(), string(*text))
		fmt.Println()
		return subcommands.ExitFailure
	}

	for _, reservation := range reservations {
		for _, i := range reservation.Instances {
			storePublicDnsNameToFile(p, i)
		}
	}

	return subcommands.ExitSuccess
}

func getText(args []string) (*string, error) {
	if len(args) == 0 {
		return nil, errors.New("search text is required")
	}

	if len(args) > 1 {
		return nil, errors.New("search text should be single")
	}

	return &args[0], nil
}

func getReservations(text *string, flags *tagCmd) ([]*ec2.Reservation, error) {

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String(flags.region),
		Credentials: credentials.NewSharedCredentials(flags.credential.filename, flags.credential.profile)})

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String(fmt.Sprintf("tag:%s", flags.tagKey)),
				Values: []*string{
					aws.String(*text),
				},
			},
		},
	}

	res, _ := svc.DescribeInstances(params)

	if len(res.Reservations) == 0 {
		return nil, errors.New("Instance not found.")
	}

	return res.Reservations, nil
}

func storePublicDnsNameToFile(flags *tagCmd, i *ec2.Instance) {
	var tagValue string
	for _, tag := range i.Tags {
		if *tag.Key == flags.tagKey {
			tagValue = *tag.Value
			break
		}
	}

	filePath := fmt.Sprintf("./%s_%s", tagValue, *i.InstanceId)
	content := []byte(*i.PublicDnsName)
	ioutil.WriteFile(filePath, content, os.ModePerm)

	fmt.Printf("Completed saving file %s, that content is '%s'", filePath, string(*i.PublicDnsName))
	fmt.Println()
}
