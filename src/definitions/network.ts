enum currentNetworkType {
  Testnet = 'testnet',
  Mainnet = 'mainnet',
}

export const currentNetwork: currentNetworkType = currentNetworkType.Mainnet; //or 'mainnet'

// import { NetworkDataType } from '../context';
// export interface NetworkDataType {
//   name: string;
//   chainId: string;
//   rpc: string[];
//   nativeCurrency: {
//     name: string;
//     symbol: string;
//     decimals: number;
//   };
//   blockExplorer: string[];
//   imageUrlRoot: string;
// }

export const networks = {
  filecoinHyperspace: {
    name: 'Filecoin Hyperspace Testnet',
    chainId: '0xc45',
    rpc: [
      'https://api.hyperspace.node.glif.io/rpc/v1',
      'https://filecoin-hyperspace.chainstacklabs.com/rpc/v1',
      'https://rpc.ankr.com/filecoin_testnet',
    ],
    nativeCurrency: {
      name: 'tFIL',
      symbol: 'tFIL',
      decimals: 18,
    },
    blockExplorer: [
      'https://fvm.starboard.ventures/transactions/',
      'https://hyperspace.filscan.io/',
    ],
    contracts: {
      WATERLILY_CONTRACT_ADDRESS: '0xfF1f3f598036BA96752b76262aB80cE4e6965eB0',
      LILYPAD_EVENTS_CONTRACT_ADDRESS:
        '0x5aFC3aCAFd6A2cB0bbdeD5A75c8d1E361FD25863',
    },
    imageUrlRoot: `https://waterlily.cluster.world/job/3141-`,
  },
  filecoinMainnet: {
    name: 'Filecoin Mainnet',
    chainId: '0x13a',
    rpc: ['https://api.node.glif.io'],
    nativeCurrency: {
      name: 'Filecoin',
      symbol: 'FIL',
      decimals: 18,
    },
    blockExplorer: ['https://filfox.info/tx/'],
    contracts: {
      WATERLILY_CONTRACT_ADDRESS: '0xdC7612fa94F098F1d7BB40E0f4F4db8fF0bC8820',
      LILYPAD_EVENTS_CONTRACT_ADDRESS:
        '0x6a46ddE41c3f572A07527149552b4B1875B0361B',
    },
    imageUrlRoot: `https://waterlily.cluster.world/job/314-`,
  },
};

export const getNetwork = () => {
  if (typeof window === "undefined") {
    return networks.filecoinHyperspace
  }
  const urlSearchParams = new URLSearchParams((window as any).location.search);
  const params = Object.fromEntries(urlSearchParams.entries());
  let currentNetworkName: string = params.waterlilyNetwork || '';
  if(currentNetworkName == 'filecoinHyperspace') return networks.filecoinHyperspace
  return networks.filecoinMainnet
}