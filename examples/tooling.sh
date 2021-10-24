#!/bin/bash
curl --silent --location "https://github.com/lanixx-hh/runfromyaml/releases/latest/download/runfromyaml-$(uname -s)-$(uname -m).tar.gz" | tar xz
curl --silent --location --output tooling.yaml https://raw.githubusercontent.com/LANIXX-HH/runfromyaml/master/examples/tooling.yaml
exec ./runfromyaml -file tooling.yaml
