import typing

from pymongo.errors import DuplicateKeyError
from pymongo import ASCENDING, DESCENDING
from bson import ObjectId

from runehistory.app.database import DatabaseAdapter, TableAdapter
from runehistory.app.exceptions import DuplicateError

if typing.TYPE_CHECKING:
    from pymongo.database import Database
    from pymongo.collection import Collection


class MongoDatabaseAdapter(DatabaseAdapter):
    def __init__(self, db: 'Database'):
        self.db = db

    def table(self, table: str, identifier: str = None) -> 'MongoTableAdapter':
        return MongoTableAdapter(self.db[table], identifier)


class MongoTableAdapter(TableAdapter):
    def __init__(self, collection: 'Collection', identifier: str = None):
        super().__init__(identifier)
        self.collection = collection

    def _record_to_id(self, record: typing.Dict) -> typing.Dict:
        if self.identifier and self.identifier in record:
            record['_id'] = record.pop(self.identifier)
        return record

    def _record_from_id(self, record: typing.Dict) -> typing.Dict:
        if '_id' in record and isinstance(record['_id'], ObjectId):
            record['_id'] = str(record['_id'])
        if self.identifier and '_id' in record:
            record[self.identifier] = record.pop('_id')
        return record

    def _projection_from_list(self, fields: typing.List = None) -> typing.Dict:
        if not fields:
            return {}
        fields = ['_id' if field is self.identifier else field for field in
                  tuple(fields)]
        projection = {field: 1 for field in fields}
        if '_id' not in fields:
            projection['_id'] = 0
        return projection

    def insert(self, record: typing.Dict) -> typing.Dict:
        try:
            record = self._record_to_id(record)
            if record['_id'] is None:
                record.pop('_id')
            self.collection.insert_one(record)
        except DuplicateKeyError:
            raise DuplicateError('Duplicate record')
        return self._record_from_id(record)

    def find_one(self, where: typing.List = None, fields: typing.List = None)\
            -> typing.Union[typing.Dict, None]:
        parsed_where = self._parse_conditions(where)
        record = self.collection.find_one(
            parsed_where,
            projection=fields
        )
        if record is None:
            return None
        return self._record_from_id(record)

    def find(self, where: typing.List = None, fields: typing.List = None,
             limit: int = 100, offset: int = None,
             order: typing.List = None
             ) -> typing.List:
        parsed_where = self._parse_conditions(where)

        results = self.collection.find(parsed_where, fields).limit(limit)
        if offset is not None:
            results = results.skip(offset)
        if order is not None:
            updated_order = []
            for item in order:
                direction = DESCENDING if item[1] is 'desc' else ASCENDING
                updated_order.append((item[0], direction))
            results = results.sort(updated_order)
        return [self._record_from_id(record) for record in results]

    def _parse_conditions(self, conditions: typing.Union[typing.List, None],
                          statement: str = 'and') -> typing.Dict:
        if not conditions:
            return dict()
        parsed_conditions = dict()
        parsed_conditions['${}'.format(statement)] = [
            self._parse_condition(condition) for condition in conditions]

        return parsed_conditions

    def _parse_condition(self, condition: typing.Union[
        typing.List, typing.Dict]) -> typing.Dict:
        if isinstance(condition, list):
            return self._parse_condition_list(condition)
        if isinstance(condition, dict):
            return self._parse_condition_dict(condition)

    def _parse_condition_list(self, condition: typing.List) -> typing.Dict:
        key = '_id' if condition[0] == self.identifier else condition[0]
        if len(condition) is 2:
            if isinstance(condition[1], dict):
                return self._parse_condition_dict(condition[1])
            return {key: {'$eq': condition[1]}}
        if len(condition) is 3:
            if condition[1] == '=':
                return {key: {'$eq': condition[2]}}
            if condition[1] == '>':
                return {key: {'$gt': condition[2]}}
            if condition[1] == '>=':
                return {key: {'$gte': condition[2]}}
            if condition[1] == '<':
                return {key: {'$lt': condition[2]}}
            if condition[1] == '<=':
                return {key: {'$lte': condition[2]}}
        raise Exception('Unhandled condition')

    def _parse_condition_dict(self, conditions: typing.Dict) -> typing.Dict:
        parsed_conditions = {}
        for statement, sub_conditions in conditions.items():
            parsed_conditions.update(
                self._parse_conditions(sub_conditions, statement))
        return parsed_conditions
