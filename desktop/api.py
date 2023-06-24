from db import db_insert, db_select


class API:
    def register_entity(self, entity_data: dict):
        db_insert('entities', entity_data)
        return entity_data

    def get_entities(self):
        return db_select('entities')

    def create_invoices(self, invoices_data):
        print(invoices_data)
        return invoices_data

    def cancel_invoices(self, cancelings_data):
        print(cancelings_data)
        return cancelings_data
