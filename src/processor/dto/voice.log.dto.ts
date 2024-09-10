// Ortak özellikleri içeren ana sınıf
export class VoiceBaseDto {
  status: string;
  guildId: string;
  userId: string;
  username: string;
  eventTime: string;
}

// İlk DTO
export class VoiceLogDto extends VoiceBaseDto {
  channelId: string;
  channelName: string;
}

// İkinci DTO
export class VoiceMoveLogDto extends VoiceBaseDto {
  oldChannelId: string;
  oldChannelName: string;
  newChannelId: string;
  newChannelName: string;
}
