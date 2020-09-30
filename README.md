# dlibox-7688

## Setup

On a Raspberry Pi:

```
wget https://raw.githubusercontent.com/periph/bootstrap/master/setup.sh
chmod +x setup.sh
./setup.sh do_golang

# Currently assuming not in go module mode:
go get -u github.com/maruel/dlibox-7688
cp ~/go/src/github.com/maruel/dlibox-7688/on-start-7688.sh .
~/go/src/github.com/maruel/dlibox-7688/setup-7688.sh
sudo shutdown -r now
```
