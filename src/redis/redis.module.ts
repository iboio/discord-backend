import { Module, Global } from '@nestjs/common';
import * as Redis from 'ioredis';
import { RedisService } from './redis.service';
import { RedisController } from './redis.controller';

@Global()
@Module({
  providers: [
    {
      provide: 'REDIS_CLIENT_DBO',
      useFactory: () => {
        return new Redis.Redis({
          host: 'localhost',
          port: 10050,
          db: 0,
        });
      },
    },
    {
      provide: 'REDIS_CLIENT_DB1',
      useFactory: () => {
        return new Redis.Redis({
          host: 'localhost',
          port: 10050,
          db: 1,
        });
      },
    },
    {
      provide: 'REDIS_CLIENT_DB2',
      useFactory: () => {
        return new Redis.Redis({
          host: 'localhost',
          port: 10050,
          db: 2,
        });
      },
    },
    RedisService,
  ],
  exports: [
    'REDIS_CLIENT_DBO',
    'REDIS_CLIENT_DB1',
    'REDIS_CLIENT_DB2',
    RedisService,
  ],
  controllers: [RedisController],
})
export class RedisModule {}
