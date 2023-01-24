#!/bin/bash
curl --silent --location "https://github.com/lanixx-hh/runfromyaml/releases/v4.0.0/download/runfromyaml-$(uname -s)-$(uname -m).tar.gz" | tar xz
curl --silent --location --output tooling.yaml https://raw.githubusercontent.com/LANIXX-HH/runfromyaml/master/examples/tooling.yaml
