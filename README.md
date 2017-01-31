# describe-ec2

This project contains command line tool of searching ec2 instance. If you install this project, you can search ec2 instance by command line tool.

## How to search instance
The describe-ec2 client can search by the ec2 tag. There is only search option now.

## Command line options
```
$ describe-ec2 help tag
tag [-credential-profile default] [-credential-filename '~/.aws/credentials'] [-region ap-northeast-1] [-tag-key Name] '*dev*' :
  Created or updated text file
  -credential-filename string
    	optional: aws credential file name, when filename is empty, that will use '$HOME/.aws/credentials'
  -credential-profile string
    	optional: aws credential profile, default value is 'default' (default "default")
  -region string
    	optional: aws region, default value is 'ap-northeast-1' (default "ap-northeast-1")
  -tag-key string
    	target tag key, default value is 'Name' (default "Name")
```

## More useful searching of ec2 instance
If you install [peco](https://github.com/peco/peco), you can get more useful searching of ec2 instance.
This project contains zsh functions that provides incremental search of ec2 instance by using peco.

## How to install
```
$ go get github.com/nsoushi/describe-ec2/cmd/describe-ec2
$ source $GOPATH/src/github.com/nsoushi/describe-ec2/.zsh.describe_ec2
$ describe-ec2 tag '*AWS*'
$ peco-describe-ec2　// You can search instance of includes '*AWS*' keyword.
```

## Additional Reports
[This post](http://naruto-io.hatenablog.com/entry/2017/01/22/214441) reported how to build this project, sorry that supported only japanese.
