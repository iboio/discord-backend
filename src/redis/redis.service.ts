// voice log db: 0
// log channel db:0
// commandData  and instance data db:2
// bull db:15
import { Inject, Injectable } from '@nestjs/common';
import Redis from 'ioredis';

@Injectable()
export class RedisService {
  private redisClients: Map<string, Redis>;
  constructor(
    @Inject('REDIS_CLIENT_DBO') private readonly redisClientDB0: Redis,
    @Inject('REDIS_CLIENT_DB1') private readonly redisClientDB1: Redis,
    @Inject('REDIS_CLIENT_DB2') private readonly redisClientDB2: Redis,
  ) {
    this.redisClients = new Map([
      ['db0', this.redisClientDB0],
      ['db1', this.redisClientDB1],
      ['db2', this.redisClientDB2],
    ]);
  }

  private getClient(dbName: string): Redis {
    const client = this.redisClients.get(dbName);
    if (!client) {
      throw new Error(`No Redis client found for database: ${dbName}`);
    }
    return client;
  }
  // redis string methods -- start
  async stringSet(db: string, string: string, key: string) {
    const client = this.getClient(db);
    await client.set(string, key);
  }
  async stringGet(db: string, string: string) {
    const client = this.getClient(db);
    return client.get(string);
  }
  async stringDel(db: string, string: string) {
    const client = this.getClient(db);
    await client.del(string);
  }
  // redis string methods -- end
  //----------------------------
  // redis list methods -- start
  async listSet(db: string, list: string, key: string) {
    const client = this.getClient(db);
    await client.lpush(list, key);
  }
  async listGet(db: string, list: string): Promise<string[]> {
    const client = this.getClient(db);
    return client.lrange(list, 0, -1);
  }
  async listExpire(db: string, list: string, second: number) {
    const client = this.getClient(db);
    await client.expire(list, second);
  }
  async listDelValue(db: string, list: string, value: string) {
    const client = this.getClient(db);
    await client.lrem(list, 0, value);
  }
  async moveValueBetweenLists(
    db: string,
    sourceList: string,
    destinationList: string,
    value: string,
  ) {
    await this.listDelValue(db, sourceList, value);
    await this.listSet(db, destinationList, value);
  }
  // redis list methods -- end
  //----------------------------
  // redis hash methods -- start
  async hashSet(db: string, hash: string, key: string, value: string) {
    const client = this.getClient(db);
    await client.hset(hash, key, value);
  }

  async hashGet(db: string, hash: string, key: string) {
    const client = this.getClient(db);
    return client.hget(hash, key);
  }

  async hashDelKey(db: string, hash: string, key: string) {
    const client = this.getClient(db);
    await client.hdel(hash, key);
  }

  async hashDel(db: string, hash: string) {
    const client = this.getClient(db);
    await client.del(hash);
  }

  async hashGetAll(db: string, hash: string) {
    const client = this.getClient(db);
    return client.hgetall(hash);
  }
  // redis hash methods -- end
  //----------------------------
  //extra
  async getAllLists(db: string): Promise<{ [key: string]: string[] }> {
    const client = this.getClient(db);
    const lists = {};
    const stream = client.scanStream({
      match: '*',
    });

    for await (const keys of stream) {
      if (keys.length) {
        const pipeline = client.pipeline();
        keys.forEach((key) => {
          pipeline.type(key);
        });

        const results = await pipeline.exec();
        const listKeys = keys.filter(
          (_, index) => results[index][1] === 'list',
        );

        if (listKeys.length) {
          const listPipeline = client.pipeline();
          listKeys.forEach((listKey) => {
            listPipeline.lrange(listKey, 0, -1);
          });

          const listResults = await listPipeline.exec();
          listResults.forEach((result, index) => {
            lists[listKeys[index]] = result[1];
          });
        }
      }
    }
    return lists;
  }
}
