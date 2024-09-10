import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { RedisModule } from './redis/redis.module';
import { BullModule } from '@nestjs/bullmq';
import { ProcessorModule } from './processor/processor.module';
import { DatabaseModule } from './database/database.module';
import { ScheduleModule } from "@nestjs/schedule";
import { CronModule } from "./cron/cron.module";
import { ConfigModule, ConfigService } from "@nestjs/config";
import { DiscordModule } from "@discord-nestjs/core";
import { GatewayIntentBits } from 'discord.js';

@Module({
  imports: [
    BullModule.forRoot({
      connection: {
        host: 'localhost',
        port: 10050,
        db: 15,
      },
    }),
    BullModule.registerQueue({
      name: 'messageLog',
    }),
    BullModule.registerQueue({
      name: 'voiceLog',
    }),
    BullModule.registerQueue({
      name: 'voiceStateLog',
    }),
    DiscordModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => ({
        token: configService.get('DISCORD_BOT_TOKEN'),
        discordClientOptions: {
          intents: [
            GatewayIntentBits.Guilds,
            GatewayIntentBits.GuildMessages,
            GatewayIntentBits.MessageContent,
            GatewayIntentBits.GuildVoiceStates,
            GatewayIntentBits.GuildModeration,
            GatewayIntentBits.GuildMembers,
            GatewayIntentBits.GuildIntegrations,
          ],
        },
      }),
    }),
    RedisModule,
    ProcessorModule,
    DatabaseModule,
    ScheduleModule.forRoot(),
    ConfigModule.forRoot(),
    CronModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
