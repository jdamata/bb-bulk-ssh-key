# bb-bulk-ssh-key
Upload multiple ssh keys to your bitbucket account. This utility will check a directory for public ssh keypairs, validate the public key and upload it to your bitbucket account.

## Installation

You can grab a pre-compiled version of bb-bulk-ssh-key in the release tab or generate your own:

```console
go get -u github.com/jdamata/bb-bulk-ssh-key
```

## Usage

### Create an app password

To create an app password:
- From your avatar in the bottom left, click Personal settings.
- Click App passwords under Access management.
- Click Create app password.
- Give the app password a name related to the application that will use the password.
- Grant write access to account
- Copy the generated password and either record or paste it into the application you want to give access. The password is only displayed this one time.

### Run 

```console
./bb-bulk-ssh-key -u joeldamata -p <REDACTED> -d ~/joeldamata/.ssh/
```
