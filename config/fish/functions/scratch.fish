function scratch --description 'Open a dated scratch directory'
    set -l d $HOME/Code/.scratch/(date +%Y-%m-%d)
    mkdir -p $d; and cd $d
end
