# s
Secure and anonymous messaging in just one letter

## Instalation

You need to have Go installed: https://github.com/golang/go

Then run:

`go get -v -d github.com/kenan-rhoton/s`

## Usage

`s srv` to start a server (for now fixed listening on port 8090)

`s gen` Generate a secure Private/Public key pair in current directory

`s sel <server>` select a server to use as your target

`s reg <alias>` register on the selected server with that alias

`s` ask the server to give you all messages sent to you

`s <user> <message> send a message to a certain user`

## Limitations

- Your identity is saved on the current directory (will move this to always be on $HOME or $S\_DIR)
- No ability to choose port
- Need more checks when retrieving messages on server to ensure no eavesdropping (but everything is still encrypted, this would be an added layer of security)

