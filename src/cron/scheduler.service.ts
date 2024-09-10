import { Injectable } from '@nestjs/common';
import * as moment from 'moment';
import { Cron, CronExpression } from '@nestjs/schedule';
import { RedisService } from '../redis/redis.service';
import { VoiceLog } from '../database/log-service/discord/voice.log';
import { DateToEpochConvertor } from '../methods/date.to.epoch.convertor';
@Injectable()
export class SchedulerService {
  constructor(
    private readonly redisService: RedisService,
    private readonly dcVoiceLogService: VoiceLog,
  ) {}
  @Cron('* * * * *')
  async voiceStatistic() {
    const eventTime = moment(new Date()).format('YYYY-MM-DD HH:mm:ss');
    await this.redisService.stringSet('db2', 'cronTime', eventTime);
    const epochTime = Math.floor(Date.now() / 1000);
    const channelList = await this.redisService.getAllLists('db0');
    for (const [key, value] of Object.entries(channelList)) {
      for (const user of value) {
        const hashData = await this.redisService.hashGetAll('db0', user);
        console.log(user, hashData);
        const eventDateTime = await DateToEpochConvertor(hashData.eventTime);
        let count = 60;
        if (epochTime - eventDateTime < 60) {
          count = epochTime - eventDateTime;
        }
        await this.dcVoiceLogService.insertVoiceLog({
          guildId: hashData.guildId,
          channelId: key,
          channelName: hashData.channelName,
          userId: user,
          username: hashData.username,
          count: count,
          eventTime: eventTime,
        });
      }
    }
  }
}
