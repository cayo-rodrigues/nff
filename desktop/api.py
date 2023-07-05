import db


class API:
    def register_entity(self, entity_data: dict):
        db.insert('entities', entity_data)

    def get_entities(self):
        return db.select('entities')

    def delete_entity(self, entity_id: int):
        db.delete('entities', entity_id)

    def update_entity(self, entity_data: dict, entity_id: int):
        db.update('entities', entity_data, entity_id)

    def create_invoices(self, invoices_data):
        print(invoices_data)
        return invoices_data

    def cancel_invoices(self, cancelings_data):
        print(cancelings_data)
        return cancelings_data

    def get_lists_data(self, columns: str):
        return db.select('lists_data', columns=columns, group_by_columns=True)
