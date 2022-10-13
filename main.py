import sys
import src.basic as basic

try:
    file_name = sys.argv[1]
    result, error = basic.run("<stdin>", f'RUN("{file_name}")')
    if error:
        print(error.as_string())
    elif result:
        if len(result.elements) == 1:
            print(repr(result.elements[0]))
        else:
            print(repr(result))
    print()
except IndexError:
    while True:
        try:
            text = input(">>> ")
            if text == "exit()":
                break
            if text.strip() == "":
                continue
            result, error = basic.run("<stdin>", text)
            if error:
                print(error.as_string())
            elif result:
                if len(result.elements) == 1:
                    print(repr(result.elements[0]))
                else:
                    print(repr(result))
            print()
        except Exception as e:
            print(f"Exception: {e}")
        except KeyboardInterrupt:
            break
