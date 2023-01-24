#!/bin/bash
curl --silent --location "https://github.com/LANIXX-HH/runfromyaml/releases/download/v4.0.0/runfromyaml-$(uname -s)-$(uname -m).tar.gz" | tar xz
curl --silent --location --output tooling.yaml https://raw.githubusercontent.com/LANIXX-HH/runfromyaml/master/examples/tooling.yaml
