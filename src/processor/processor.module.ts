import { Module } from '@nestjs/common';
import { MessageProcessor } from './discord/message.processor';
import { DatabaseModule } from '../database/database.module';
import { VoiceProcessor } from './discord/voice.processor';

@Module({
  imports: [DatabaseModule],
  providers: [MessageProcessor, VoiceProcessor],
  controllers: [],
  exports: [],
})
export class ProcessorModule {}
