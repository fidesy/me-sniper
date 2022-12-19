# me-sniper

me-sniper is an on-chain MagicEden sniper that detects all new listings on MagicEden throw blockchain with telegram notifications feature.

## Installation

1. Clone repository.
    ```bash
    git clone https://github.com/fidesy/me-sniper.git
     ```
2. Create .env file with two variables:

    ---
    NODE_ENDPOINT=<YOUR_SOLANA_RPC_ENDPOINT>

    TELEGRAM_APIKEY=<YOUR_TELEGRAM_APIKEY>
    
    ---
    If you don't need telegram notification then just skip second variable.
    
    You can get Solana RPC node for free at https://www.quicknode.com
3. Run script
    ```bash
    go run cmd/main.go
    ```
4. (Optional) If you are using telegram bot, then write /start command to it.
## Usage

### Logs contain:

* Action type (now only list/buy)
* Block/Current timestamp
* collection symbol, token name
* price, rarity, rank, supply, seller, buyer

![](./data/logs.png)


### Telegram notifications.

![](./data/telegram.png)

## License
[MIT](https://choosealicense.com/licenses/mit/)