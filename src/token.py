from src.position import Position


class Token:
    def __init__(self, type_: str, value: str=None, pos_start: Position=None, pos_end: Position=None):
        self.type: str = type_
        self.value: str = value

        if pos_start:
            self.pos_start = pos_start.copy()
            self.pos_end = pos_start.copy()
            self.pos_end.advance()

        if pos_end:
            self.pos_end = pos_end.copy()

    def matches(self, type_: str, value: str):
        return self.type == type_ and self.value == value

    def __repr__(self):
        if self.value:
            return f"{self.type}:{self.value}"
        return f"{self.type}"
