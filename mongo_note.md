

  1.  mongodump -d doc_manager -o doc.dmp
  2. scp -r doc.dmp root@118.89.108.25:/root/go-eladmin/db/
  3. ssh root@...
  4. /usr/local/mongodb/bin/mongorestore -u doc -p doc11121014a -d docmanager doc.dmp/doc_manager




