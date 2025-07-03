use anyhow::Result;
use redis::aio::ConnectionManager;
use serde::{Deserialize, Serialize};
use std::time::Duration;
use tracing::{error, info, instrument};

#[derive(Clone)]
pub struct CacheService {
    redis: ConnectionManager,
}

impl CacheService {
    pub async fn new(redis: ConnectionManager) -> Result<Self> {
        info!("Initializing cache service");
        Ok(Self { redis })
    }

    #[instrument(skip(self))]
    pub async fn ping(&self) -> Result<()> {
        let mut conn = self.redis.clone();
        let _: String = redis::cmd("PING").query_async(&mut conn).await?;
        Ok(())
    }

    #[instrument(skip(self, value))]
    pub async fn set<T>(&self, key: &str, value: &T, ttl: Option<Duration>) -> Result<()>
    where
        T: Serialize,
    {
        let mut conn = self.redis.clone();
        let serialized = serde_json::to_string(value)?;
        
        if let Some(ttl) = ttl {
            redis::cmd("SETEX")
                .arg(key)
                .arg(ttl.as_secs())
                .arg(serialized)
                .query_async(&mut conn)
                .await?;
        } else {
            redis::cmd("SET")
                .arg(key)
                .arg(serialized)
                .query_async(&mut conn)
                .await?;
        }
        
        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn get<T>(&self, key: &str) -> Result<Option<T>>
    where
        T: for<'de> Deserialize<'de>,
    {
        let mut conn = self.redis.clone();
        let result: Option<String> = redis::cmd("GET")
            .arg(key)
            .query_async(&mut conn)
            .await?;
            
        match result {
            Some(data) => {
                let deserialized = serde_json::from_str(&data)?;
                Ok(Some(deserialized))
            }
            None => Ok(None),
        }
    }

    #[instrument(skip(self))]
    pub async fn delete(&self, key: &str) -> Result<bool> {
        let mut conn = self.redis.clone();
        let result: i32 = redis::cmd("DEL")
            .arg(key)
            .query_async(&mut conn)
            .await?;
        Ok(result > 0)
    }

    #[instrument(skip(self))]
    pub async fn exists(&self, key: &str) -> Result<bool> {
        let mut conn = self.redis.clone();
        let result: i32 = redis::cmd("EXISTS")
            .arg(key)
            .query_async(&mut conn)
            .await?;
        Ok(result > 0)
    }

    #[instrument(skip(self))]
    pub async fn increment(&self, key: &str, value: i64) -> Result<i64> {
        let mut conn = self.redis.clone();
        let result: i64 = redis::cmd("INCRBY")
            .arg(key)
            .arg(value)
            .query_async(&mut conn)
            .await?;
        Ok(result)
    }

    #[instrument(skip(self))]
    pub async fn expire(&self, key: &str, ttl: Duration) -> Result<bool> {
        let mut conn = self.redis.clone();
        let result: i32 = redis::cmd("EXPIRE")
            .arg(key)
            .arg(ttl.as_secs())
            .query_async(&mut conn)
            .await?;
        Ok(result > 0)
    }
}
