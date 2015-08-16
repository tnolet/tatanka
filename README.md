# Tatanka
Tatanka is a migrating and daemon that lives of the global AWS spot instance infrastructure. Tatanka kills its current self after a specified amount of time and later reincarnates somewhere else. Tatanka communicates by email and Twitter. 


## Commands

run tatanka in test mode

```
./tatanka -localConfig="conf/local_test.json" -noop=true
```
    
## TODO
[ ] Allow simple or advanced price strategy 
[ ] Retry cancelling of spot requests when request not found.  
[ ] AWS Dry run with NOOP mode.  
[ ] allow lifetime of 0 or less for indefinite lifetime
[x] Do not make an extra spot requests on evac when the initial spot request has already been honored.  
[x] validate time ranges input (e.g. bid offset < lifetime).  
