import { Processor, WorkerHost } from '@nestjs/bullmq';
import { Job } from 'bullmq';
import {
  VoiceBaseDto,
  VoiceLogDto,
  VoiceMoveLogDto,
} from '../dto/voice.log.dto';
import { RedisService } from '../../redis/redis.service';
import { DateToEpochConvertor } from '../../methods/date.to.epoch.convertor';
import { VoiceLog } from '../../database/log-service/discord/voice.log';
import * as moment from 'moment';

@Processor('voiceLog')
export class VoiceProcessor extends WorkerHost {
  constructor(
    private readonly redisService: RedisService,
    private readonly dcVoiceLogService: VoiceLog,
  ) {
    super();
  }

  async process(job: Job) {
    await this.voiceStatusHandler(job.data);
  }

  async voiceStatusHandler(data: VoiceBaseDto) {
    if (data.status == 'join') {
      await this.userLogin(data as VoiceLogDto);
    } else if (data.status == 'left') {
      await this.userLeft(data as VoiceLogDto);
    } else if (data.status == 'move') {
      await this.userMove(data as VoiceMoveLogDto);
    }
  }

  async userLogin(data: VoiceLogDto) {
    console.log('user join', data.eventTime);
    await this.redisService.listSet('db0', data.channelId, data.userId);
    const { guildId, channelName, username, eventTime } = data;
    const newData = { guildId, channelName, username, eventTime };
    for (const key in newData) {
      await this.redisService.hashSet('db0', data.userId, key, newData[key]);
    }
  }

  async userLeft(data: VoiceLogDto) {
    await this.whenUserLeftInsertVoiceCount(
      data,
      data.channelName,
      data.channelId,
      await this.redisService.hashGet('db0', data.userId, 'eventTime'),
    );
    console.log('user left', data.eventTime);
    await this.redisService.listDelValue('db0', data.channelId, data.userId);
    await this.redisService.hashDel('db0', data.userId);
  }

  async userMove(data: VoiceMoveLogDto) {
    console.log('user move', data.eventTime);
    await this.whenUserLeftInsertVoiceCount(
      data,
      data.oldChannelName,
      data.oldChannelId,
      await this.redisService.hashGet('db0', data.userId, 'eventTime'),
    );
    await this.redisService.moveValueBetweenLists(
      'db0',
      data.oldChannelId,
      data.newChannelId,
      data.userId,
    );
    await this.redisService.hashDelKey('db0', data.userId, 'channelName');
    await this.redisService.hashSet(
      'db0',
      data.userId,
      'channelName',
      data.newChannelName,
    );
    await this.redisService.hashSet(
      'db0',
      data.userId,
      'eventTime',
      data.eventTime,
    );
  }

  async whenUserLeftInsertVoiceCount(
    data: VoiceBaseDto,
    channelName: string,
    channelId: string,
    time: any,
  ) {
    const cronTime = await DateToEpochConvertor(
      await this.redisService.stringGet('db2', 'cronTime'),
    );
    const eventTime = await DateToEpochConvertor(time);
    const epochTime = Math.floor(Date.now() / 1000);
    let count = epochTime - cronTime;
    console.log(epochTime - eventTime);
    if (cronTime < eventTime && cronTime + 60 > epochTime) {
      count = epochTime - eventTime;
      console.log(count);
    }
    await this.dcVoiceLogService.insertVoiceLog({
      guildId: data.guildId,
      channelId: channelId,
      channelName: channelName,
      userId: data.userId,
      username: data.username,
      count: count,
      eventTime: moment(new Date()).format('YYYY-MM-DD HH:mm:ss'),
    });
  }
}
