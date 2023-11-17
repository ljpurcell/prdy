# prdy
A command line tool to check for debug code (`console.log()` for JavaScript, `dd()` for Php, etc.) and run the test suite before submitting a pull-request.

Images and details for `version-1` are below. The first version was used to explore the problem space and work out any major pain points.

The second version is currently under construction and makes use popular of libraries to enhance functionality and user-experience.

### Usage
If the command is run, and a configuration file is not found, the program will automatically enter the config menu.
<img width="1172" alt="Screenshot 2023-11-17 at 5 01 48 pm" src="https://github.com/ljpurcell/prdy/assets/65317064/426dac12-de9e-49f6-9290-d0aed9217faa">

### Output
Prints the file name, line number and line contents of any lines that qualify.
<img width="1312" alt="Screenshot 2023-11-17 at 5 02 23 pm" src="https://github.com/ljpurcell/prdy/assets/65317064/82d73046-d2ab-488c-afa3-d5f112c97aa1">

### Options & Flags
#### Version 1
Currently, only supports a single `-c` flag to enter the configuration menu.
