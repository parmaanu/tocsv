# Things to do:

## TODO, rename the package to filterlogs

1. Modify config structure to allow grouping for different use cases or apps. 
    - filterlogs should be smart enough to identity the app only just through logs.
    - there should strictly be no extra column
    - read config from multiple yaml files [DONE]

2. Support N/F and N/A
    - A field is N/F if its tag is given but not found in the line
    - A field is N/A if it is not applicable for that record
    - So, if there is N/F then there is something wrong with the config

3. Support log blocks
    - Each log block will go into a single csv record. 
    - This would require specifying StartBlockPattern and EndBlockPattern

4. Easy and automated config generation using some TUI or interactive menu

5. Support switching into the loglines and the output csv or go to corresponding log block
    - It should be pretty fast to switch in and switch out
    - We can also show a height configured mini logs display at the bottom of the screen
    - To make it efficient, we should be able to catch the data fast and should be able to read the huge logs in both the directions

6. Support a live mode
    - Display csv records as logs appear in real time. 
    - This is like a real time dashboard that will monitor a file and generate the dashboard

## Small features
- Read config from different yaml files and combine them. There is a lot of error handling required here
- Do not display a column if that field is empty
- Make filterlogs a library which can be imported into another projects for data parsing and keys could be in the config so that the app can use it from the map returned by filterlogs
- Show logLines and their patterns
