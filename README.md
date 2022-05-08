# mailchimp-data-importer

The goal of this project is to sync data between two different API's, performing one GET and sending a POST.

## Installation

Please complete the config on the [config.yaml](./config/config.yaml) with the API's information. Then run 

```bash
go mod tidy
```
To install the dependencies.

Other useful commands, such as unit tests, are on the [Makefile](./Makefile).

## Usage

The project is not Dockerized (yet!). So to run the app you should run the command (with config already filled):

```bash
make run
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)