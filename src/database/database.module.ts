import { Module } from '@nestjs/common';
import { MessageLog } from './log-service/discord/message.log';
import { VoiceLog } from './log-service/discord/voice.log';
import { DummyDataService } from './dummy.data.service';
import { DatabaseController } from './database.controller';
import { DiscordModule } from '@discord-nestjs/core';
import { DiscordGuild } from './log-service/discord/discord.guild';

@Module({
  imports: [DiscordModule.forFeature()],
  providers: [MessageLog, VoiceLog, DummyDataService, DiscordGuild],
  controllers: [DatabaseController],
  exports: [MessageLog, VoiceLog],
})
export class DatabaseModule {}
