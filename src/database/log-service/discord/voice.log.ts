import { Injectable } from '@nestjs/common';
import { ClickHouseClient, createClient } from '@clickhouse/client';

@Injectable()
export class VoiceLog {
  private clickhouse: ClickHouseClient;
  constructor() {
    this.clickhouse = createClient({
      url: 'http://localhost:10100',
      database: 'discord',
    });
  }
  async insertVoiceLog(data: any) {
    await this.clickhouse.insert({
      table: 'voice',
      format: 'JSONEachRow',
      values: [data],
    });
  }
}
