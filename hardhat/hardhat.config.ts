import { HardhatUserConfig, task } from 'hardhat/config';
import '@nomicfoundation/hardhat-toolbox';
import '@openzeppelin/hardhat-upgrades';
import 'hardhat-deploy';
import { config as dotenvConfig } from 'dotenv';
import { resolve } from 'path';

const dotenvConfigPath: string = process.env.DOTENV_CONFIG_PATH || './.env';
dotenvConfig({ path: resolve(__dirname, dotenvConfigPath) });

const walletPrivateKey: string | undefined = process.env.WALLET_PRIVATE_KEY;
if (!walletPrivateKey) {
  throw new Error('Please set your Wallet private key in a .env file');
}

const config: HardhatUserConfig = {
  solidity: '0.8.17',
  defaultNetwork: 'hardhat',
  namedAccounts: {
    admin: 0,
  },
  networks: {
    hardhat: {},
    localhost: {},
    filecoinHyperspace: {
      url: 'https://api.hyperspace.node.glif.io/rpc/v1', //'https://rpc.ankr.com/filecoin_testnet', //https://filecoin-hyperspace.chainstacklabs.com/rpc/v1, "https://hyperspace.filfox.info/rpc/v1", wss://wss.hyperspace.node.glif.io/apigw/lotus/rpc/v1
      chainId: 3141,
      accounts: [walletPrivateKey],
    },
    filecoinMainnet: {
      url: 'https://api.node.glif.io', //'https://rpc.ankr.com/filecoin_testnet', //https://filecoin-hyperspace.chainstacklabs.com/rpc/v1
      chainId: 314,
      accounts: [walletPrivateKey],
    },
  },
};

export default config;
