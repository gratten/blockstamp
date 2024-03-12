#!/usr/bin/env bash

# install nix-direnv
nix profile install nixpkgs#nix-direnv

# mkdir -p $HOME/.config/direnv
# echo "source $HOME/.nix-profile/share/nix-direnv/direnvrc" >> $HOME/.config/direnv/direnvrc

# Check if the shell is bash
if [ -n "$BASH_VERSION" ]; then
    echo 'eval "$(direnv hook $SHELL)"' >> $HOME/.bashrc
    echo "Added direnv hook to ~/.bashrc"
    source $HOME/.bashrc
# Check if the shell is zsh
elif [ -n "$ZSH_VERSION" ]; then
    echo 'eval "$(direnv hook $SHELL)"' >> $HOME/.zshrc
    echo "Added direnv hook to ~/.zshrc"
    source $HOME/.zshrc
else
    echo "Unsupported shell; this script is only configured for bash or zsh"
    exit 1
fi

