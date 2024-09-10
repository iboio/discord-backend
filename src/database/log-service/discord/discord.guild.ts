import { Injectable } from '@nestjs/common';
import { ClickHouseClient, createClient } from '@clickhouse/client';
import { InjectDiscordClient } from '@discord-nestjs/core';
import { Client, Guild } from 'discord.js';

@Injectable()
export class DiscordGuild {
  private clickhouse: ClickHouseClient;

  constructor(
    @InjectDiscordClient()
    private readonly client: Client,
  ) {
    this.clickhouse = createClient({
      url: 'http://localhost:10100',
      database: 'discord',
    });
  }
  async getGuildList() {
    const guild = this.client.guilds.cache;
    return guild.map((guild: Guild) => ({
      guildId: guild.id,
      name: guild.name,
      icon: guild.iconURL(),
    }));
  }
}
