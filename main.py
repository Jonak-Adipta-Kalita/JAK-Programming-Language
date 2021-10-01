import src.basic as basic

while True:
    try:
        text = input("JAK_Programming_Language>>> ")
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
    except KeyboardInterrupt as e:
        break
