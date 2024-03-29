import { FC, useContext, useState, useEffect } from 'react';
import { Button } from '@mui/material';
import { WalletContext } from '@/context';

const style = {
  backgroundSize: '200% 200%',
  animation: `gradient-animation 4s ease infinite`,
  color: 'white',
  fontWeight: 'bold',
  fontSize: 'larger',
  paddingTop: '0.5em',
  paddingRight: '1.5em',
  paddingLeft: '1.5em',
};

export const WalletButton: FC = () => {
  const buttonStates: any = {
    connectState: {
      name: 'Connect',
      action: () => connectWallet(),
      background: `-webkit-linear-gradient(left, #f53ebb 10%, #b583ff 70%);`,
    },
    connectedState: {
      name: 'Connected',
      action: () => disconnectWallet(),
      background: `-webkit-linear-gradient(left, #30ccff 10%, #0055ff 70%);`,
    },
    installWalletState: {
      name: 'Install Wallet',
      action: () =>
        window.open(
          'https://metamask.io/download.html',
          '_blank',
          'noreferrer'
        ),
      // background: `-webkit-linear-gradient(left, #f8a929 10%, #f38218 70%);`,
    },
  };
  const { walletState, connectWallet, disconnectWallet } =
    useContext(WalletContext);

  if (!walletState) return <></>;

  const buttonState = !walletState?.web3
    ? buttonStates.installWalletState
    : walletState?.isConnected
    ? buttonStates.connectedState
    : buttonStates.connectState;

  return (
    <Button
      variant={walletState.web3 ? 'contained' : 'outlined'}
      onClick={buttonState?.action || console.log('inactive')}
      sx={{ ...style, background: `${buttonState?.background}` }}
    >
      {buttonState?.name || 'Connect'}
    </Button>
  );
};
