import sqlite3


def db_connection(f):
    def wrapper(*args, **kwargs):
        connection = sqlite3.connect('user_data.db')
        cursor = connection.cursor()
        cursor.row_factory = sqlite3.Row

        result = f(*args, **kwargs, cursor=cursor)

        connection.commit()
        connection.close()

        return [dict(row) for row in result] if result else None

    return wrapper


@db_connection
def insert(table_name: str, data: dict, cursor: sqlite3.Cursor = None):
    columns = ', '.join(data.keys())
    values = list(data.values())
    placeholders = ', '.join('?' * len(values))

    cursor.execute(
        f"INSERT INTO {table_name} ({columns}) VALUES ({placeholders})",
        values
    )


@db_connection
def select(table_name: str, cursor: sqlite3.Cursor = None):
    return cursor.execute(f"SELECT * FROM {table_name}").fetchall()


@db_connection
def delete(table_name: str, row_id: int, cursor: sqlite3.Cursor = None):
    cursor.execute(f"DELETE FROM {table_name} WHERE id = ?", [row_id])
