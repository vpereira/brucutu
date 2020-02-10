Brucutu - The bruteforce tool with the coolest test suite

Old problem, new tool. Having many times had issues with Hydra, I thought it would be a good idea to actually stand on the shoulders of giants, and use battle tested protocol decoders instead of custom decoders as Hydra does. Beside it, I would like to have a better test infrastructure and native support integrations with tools that I use daily.

![brucutu](brucutu.jpg)


Usage: ./brucutu -h

The flags used to be compatible with hydra, but I'm kind of give up on that.


TODO:

    - Better restore option
    - Heuristic switching among protocols to evade monitoring systems
    - parallelism 
    - quiet mode
    - add more protocols
    - connect to graphite and maybe grafana to better monitoring


Status:

[![CircleCI](https://circleci.com/gh/vpereira/brucutu.svg?style=svg&circle-token=ac317a178e248d31fba8efd6352af94acada1f5b)](https://circleci.com/gh/vpereira/brucutu)
