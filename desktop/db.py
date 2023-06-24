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
def db_insert(table_name: str, data: dict, cursor: sqlite3.Cursor = None):
    columns = ', '.join(data.keys())
    values = list(data.values())
    placeholders = ', '.join('?' * len(values))

    cursor.execute(
        f"INSERT INTO {table_name} ({columns}) VALUES ({placeholders})",
        values
    )


@db_connection
def db_select(table_name: str, cursor: sqlite3.Cursor = None):
    return cursor.execute(f"SELECT * FROM {table_name}").fetchall()
