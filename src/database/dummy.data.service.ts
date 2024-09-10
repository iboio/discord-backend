import { Injectable, OnModuleInit, OnModuleDestroy } from '@nestjs/common';
import { ClickHouseClient, createClient } from '@clickhouse/client';
import * as moment from 'moment';

@Injectable()
export class DummyDataService implements OnModuleInit, OnModuleDestroy {
  private clickhouse: ClickHouseClient;
  private batch = [];
  private readonly BATCH_SIZE = 1000;
  private intervalId: NodeJS.Timeout;

  constructor() {
    this.clickhouse = createClient({
      url: 'http://localhost:10100',
      database: 'discord',
    });
  }

  async getRandomElementFromObject(obj: Record<string, any>): Promise<any> {
    const entries = Object.entries(obj);
    const randomIndex = Math.floor(Math.random() * entries.length);
    const [key, value] = entries[randomIndex];
    return { key, value };
  }

  private sanitizeString(str: string): string {
    return str
      .trim()
      .replace(/[\r\n]+/g, ' ')
      .replace(/'/g, "''");
  }

  private async insertBatch(): Promise<void> {
    try {
      await this.clickhouse.insert({
        table: 'message',
        values: this.batch,
        format: 'JSONEachRow',
      });
      this.batch = []; // Clear batch after successful insertion
    } catch (error) {
      console.error('Error inserting batch:', error);
    }
  }

  private addToBatch(data: any) {
    this.batch.push(data);
    if (this.batch.length >= this.BATCH_SIZE) {
      this.insertBatch().then(() => {});
    }
  }

  async dataGenerator() {
    const guilds = {
      '1043928530012086283': {
        chat: '1043928829112090735',
        secondChat: '1126140142785134674',
        general: '1044295172588060865',
        post: '1275142432580960299',
        event: '1044295172588060865',
      },
    };
    const users = {
      mortatilki: '1001240762312314980',
      lordmusamba: '335794127361015809',
      lekodias: '365912495174451200',
      siniftakaldim: '489191368212742164',
      alperage: '551398690812854294',
      fenrir: '688580719232286879',
    };

    const { key: username, value: userId } =
      await this.getRandomElementFromObject(users);
    const { key: guildId, value: channels } =
      await this.getRandomElementFromObject(guilds);
    const { key: channelName, value: channelId } =
      await this.getRandomElementFromObject(channels);
    return {
      guildId: guildId,
      messageId: BigInt(Math.floor(Math.random() * 9e17) + 1e17).toString(), // 18 haneli rastgele sayı
      channelId: channelId,
      channelName: channelName,
      userId: userId,
      username: username,
      eventTime: moment().subtract(Math.floor(Math.random() * 7), 'days').format('YYYY-MM-DD HH:mm:ss'),
    };
  }

  onModuleInit() {
    // this.intervalId = setInterval(async () => {
    //   const data = await this.dataGenerator();
    //   this.addToBatch(data);
    // }, 100);
  }

  onModuleDestroy() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }
}
