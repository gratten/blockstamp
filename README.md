# blockstamp: Convert Any Date To The Corresponding Bitcoin Block Height

## Getting Started

### Install Nix package manager

- Use this command to install the Nix package manager on your Unix-like system (Linux/MacOS/WSL):
```shell
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install
````
You can read more about the installer [here](https://zero-to-nix.com/concepts/nix-installer).

### Activate Development Shell

- In the project directory, use this command:
```shell
nix develop
```
The dependencies will be downloaded/installed in a contained environment, and your development shell will be loaded.


### Convenience Feature: Direnv (Optional But Recommended)

You can make development more convenient by installing `nix-direnv`.  Run the setup script:
```shell
./setup-direnv.sh
```

After installing, then the first time you navigate into the project directory, you will need to allow direnv to run with the command:
```shell
direnv allow
```

Now, you don't need to type `nix develop` anymore.  Any time you change into the project directory, direnv will automatically load the development shell, and it will automatically unload the development shell when you navigate out of the project directory.  An additional benefit is that it will integrate with your existing shell of choice, e.g. bash, zsh, etc.

### Configure .env file

First create a `.env` file by copying the `.env.sample` file:
```
cp .env.sample .env
```env 

Then, in the `.env` file, replace `my_bitcoind_username` and `my_bitcoind_password` with the values from your configured Bitcoin Node.
If your Bitcoin Node is running on your local environment, you can leave "host" as localhost:8332 in the .env file.
Otherwise, you can change localhost to the IP address of your node.

### Run The Program
```bash
go run main.go
```
Then open a browser to `localhost:8000`
