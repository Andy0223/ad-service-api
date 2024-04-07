db = db.getSiblingDB('2024DcardBackend');

db.createUser({
    user: 'test',
    pwd: 'test',
    roles: [{
        role: 'readWrite',
        db: '2024DcardBackend'
    }]
});