import { Processor, WorkerHost } from '@nestjs/bullmq';
import { Job } from 'bullmq';
import { MessageLog } from '../../database/log-service/discord/message.log';

@Processor('messageLog')
export class MessageProcessor extends WorkerHost {
  constructor(private readonly dcMessageService: MessageLog) {
    super();
  }
  async process(job: Job): Promise<any> {
    await this.insertData(job.data);
    return Promise.resolve(undefined);
  }

  async insertData(data: any) {
    await this.dcMessageService.insertMessageLog(data);
  }
}
