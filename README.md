## Brucutu - The bruteforce tool with the coolest test suite

### Status:

[![CI](https://github.com/vpereira/brucutu/actions/workflows/ci.yml/badge.svg)](https://github.com/vpereira/brucutu/actions/workflows/ci.yml) [![codebeat badge](https://codebeat.co/badges/50176c03-cd69-40ad-bbc7-1a28e85326c5)](https://codebeat.co/projects/github-com-vpereira-brucutu-master)
### Description
Old problem, new tool. Having many times had issues with Hydra, I thought it would be a good idea to actually stand on the shoulders of giants, and use battle tested protocol decoders instead of custom decoders as Hydra does. Beside it, I would like to have a better test infrastructure and native support integrations with tools that I use daily.

![brucutu](brucutu.jpg)


### Usage: 

./brucutu -h

The flags used to be compatible with hydra, but I'm kind of give up on that.


TODO:

    - Better restore option
    - Heuristic switching among protocols to evade monitoring systems
    - Parallelism 
    - Quiet mode
    - Add more protocols
    - Send metrics to Prometheus