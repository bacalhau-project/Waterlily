enum currentNetworkType {
  Testnet = 'testnet',
  Mainnet = 'mainnet',
}

export const currentNetwork: currentNetworkType = currentNetworkType.Mainnet; //or 'mainnet'

export const DEFAULT_API_SERVER = 'https://api.waterlily.cluster.world'

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
      'https://fvm.starboard.ventures/hyperspace/explorer/tx/',
      'https://hyperspace.filscan.io/',
    ],
    contracts: {
      WATERLILY_CONTRACT_ADDRESS: '0xd8Cc461fe24F1ACEBB7e8ffcA8955eaB5341c19f',
      WATERLILY_NFT_CONTRACT_ADDRESS:
        '0x3619c1f295B3081985e581Ea3b8546CE629C5A3D',
    },
    imageUrlRoot: `https://waterlily.cluster.world/job/3141-`,
  },
  filecoinMainnet: {
    name: 'Filecoin Mainnet',
    chainId: '0x13a',
    rpc: ['https://api.node.glif.io'],
    // wss: ['wss://wss.node.glif.io/apigw/lotus/rpc/v1'],
    nativeCurrency: {
      name: 'Filecoin',
      symbol: 'FIL',
      decimals: 18,
    },
    blockExplorer: [
      'https://fvm.starboard.ventures/explorer/tx/',
      'https://filfox.info/tx/',
    ],
    contracts: {
      WATERLILY_CONTRACT_ADDRESS: '0xdC7612fa94F098F1d7BB40E0f4F4db8fF0bC8820',
      WATERLILY_NFT_CONTRACT_ADDRESS: '',
    },
    imageUrlRoot: `https://waterlily.cluster.world/job/314-`,
  },
};

export const getNetwork = () => {
  if (typeof window === 'undefined') {
    return networks.filecoinHyperspace;
  }
  if(window.location && window.location.hostname == 'localhost') {
    return networks.filecoinHyperspace;
  }
  const urlSearchParams = new URLSearchParams((window as any).location.search);
  const params = Object.fromEntries(urlSearchParams.entries());
  let currentNetworkName: string = params.waterlilyNetwork || '';
  if (currentNetworkName == 'filecoinHyperspace')
    return networks.filecoinHyperspace;
  return networks.filecoinMainnet;
};

export const getAPIServer = (path: string = '') => {
  if (typeof window === 'undefined') {
    return DEFAULT_API_SERVER;
  }
  const urlSearchParams = new URLSearchParams((window as any).location.search);
  const params = Object.fromEntries(urlSearchParams.entries());
  const host = params.testAPI ? 'http://localhost:3500' : DEFAULT_API_SERVER;
  return `${host}/api/v1${path}`
};
