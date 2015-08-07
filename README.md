# Tatanka
Tatanka is a migrating and daemon that lives of the global AWS spot instance infrastructure. Tatanka kills its current self after a specified amount of time and later reincarnates somewhere else. Tatanka communicates by email and Twitter. 


## Commands

run tatanka in test mode

```
./tatanka -localConfig="conf/local_test.json" -noop=true
```
    
