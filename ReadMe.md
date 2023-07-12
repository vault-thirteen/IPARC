# IPARC

## IP Address Range Collection 

This is a library to work with ranges of IP addresses.

This library enables to read a database of IP addresses stored in a CSV file 
and compressed as a ZIP file into an object in memory (RAM) and make fast 
search queries to this database. 

### Details

The database contains a sequential list of IP address ranges associated with 
various data. In this example data consists of a country associated with the 
range. 

Countries are represented by standardized two-letter codes and names in English 
language. Each IP address range stores a pointer to an object, called 'struct' 
in Go programming language. Each data object may have various fields which are 
very easy to expand as the code is written in such a manner which allows very 
easy additions of new classes or types. 

### Database Source

Database is taken from the https://lite.ip2location.com website and is free to 
use.

The source file was slightly modified in following ways:  
1. Lines with empty country were standardized for the reader.
2. Double quotes were removed from there where they were not necessary.

**Last update time of the database**: 2023-07-03.

## Usage

This library is ready to be used as a search database for IP-based geolocation.

### What is IP-based Geolocation ?

IP-based geolocation is the mapping of an IP address to the real-world 
geographic location of an Internet-connected computing or a mobile device.

## Performance

Stress test shows an average RPS of about 32M for the whole test on a decent
hardware. The test iterates through all possible IPv4 addresses, i.e. from
0.0.0.0 to 255.255.255.255. Search algorithm uses binary search in array.
