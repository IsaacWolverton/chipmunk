import time

"""
 This program will continously count from 0 to 2^16 inclusively
 and output the current value every 1 seconds to stdout. This 
 program is effectively the same as simplecounter, however, 
 in a high level language.
  Parameters: none
  Return: none (unreachable during normal execution)
""" 
def main():
    print("Welcome from simplepython")

    var = 0
    while True:
        print(var)

        time.sleep(1)
        var += 1

        # Loop back to 0 if var is bigger than max uint16
        if var > 2**16:
            var = 0

"""
 Start the program if file is ran
"""
if __name__ == "__main__":
    main()