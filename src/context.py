from src.position import Position
from src.symbol_table import SymbolTable


class Context:
    def __init__(
        self,
        display_name: str,
        parent: SymbolTable = None,
        parent_entry_pos: Position = None,
    ):
        self.display_name: str = display_name
        self.parent: Context = parent
        self.parent_entry_pos: Position = parent_entry_pos
        self.symbol_table: SymbolTable = None
