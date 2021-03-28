# dlibox-7688

## Setup

On a Raspberry Pi:

```
# Install latest Go:
wget https://raw.githubusercontent.com/periph/bootstrap/master/setup.sh
bash setup.sh do_golang
rm setup.sh

# Currently assuming not in go module mode:
go get -u github.com/maruel/dlibox-7688
~/go/src/github.com/maruel/dlibox-7688/setup-7688.sh
sudo shutdown -r now
```
