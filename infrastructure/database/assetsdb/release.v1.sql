--- Release 1
----- Tables
CREATE TABLE dbo.Asset
(
    ID bigint IDENTITY(1,1) NOT NULL,
    AssetType nvarchar(255) NOT NULL,
    Name nvarchar(255) NOT NULL,
    IsActive BIT NOT NULL DEFAULT (1),
    CONSTRAINT PK_Asset PRIMARY KEY (ID),
    CONSTRAINT CHK_Asset_AssetType CHECK (AssetType IN ('Stock', 'Crypto'))
)

CREATE TABLE dbo.StockAsset
(
    ID bigint NOT NULL,
    Ticker nvarchar(255) NOT NULL,
    Exchange nvarchar(255) NOT NULL,
    Currency nvarchar(255) NOT NULL,
    CONSTRAINT PK_StockAsset PRIMARY KEY (ID),
    CONSTRAINT FK_StockAsset_Asset FOREIGN KEY (ID) REFERENCES dbo.Asset (ID)
)

CREATE TABLE dbo.CryptoAsset
(
    ID bigint NOT NULL,
    Ticker nvarchar(255) NOT NULL,
    CONSTRAINT PK_CryptoAsset PRIMARY KEY (ID),
    CONSTRAINT FK_CryptoAsset_Asset FOREIGN KEY (ID) REFERENCES dbo.Asset (ID)
)

----- Indexes
CREATE UNIQUE INDEX IX_StockAsset_Ticker ON dbo.StockAsset (Ticker)

CREATE UNIQUE INDEX IX_CryptoAsset_Ticker ON dbo.CryptoAsset (Ticker)

----- Stored Procedures
------- Stock Assets
CREATE PROCEDURE dbo.CreateStockAsset
    @Name nvarchar(255),
    @Ticker nvarchar(255),
    @Exchange nvarchar(255),
    @ISIN nvarchar(255)
AS 
BEGIN
    SET NOCOUNT ON;

    DECLARE @ID bigint

    INSERT INTO dbo.Asset (AssetType, Name)
    VALUES ('Stock', @Name)

    SET @ID = SCOPE_IDENTITY()

    INSERT INTO dbo.StockAsset (ID, Ticker, Exchange, ISIN)
    VALUES (@ID, @Ticker, @Exchange, @ISIN)

    SELECT @ID
END

CREATE PROCEDURE dbo.GetStockAssetByTicker
    @Ticker nvarchar(255)
AS
BEGIN
    SET NOCOUNT ON;

    SELECT
        A.ID,
        A.Name,
        S.Ticker,
        S.Exchange,
        S.ISIN
    FROM dbo.Asset A
    INNER JOIN dbo.StockAsset S ON S.ID = A.ID
    WHERE S.Ticker = @Ticker
END

CREATE PROCEDURE dbo.ListStockAssets
AS
BEGIN
    SET NOCOUNT ON;

    SELECT
        A.ID,
        A.Name,
        S.Ticker,
        S.Exchange,
        S.ISIN
    FROM dbo.Asset A
    INNER JOIN dbo.StockAsset S ON S.ID = A.ID
END

------- Crypto Assets
CREATE PROCEDURE dbo.CreateCryptoAsset
    @Name nvarchar(255),
    @Ticker nvarchar(255)
AS
BEGIN
    SET NOCOUNT ON;

    DECLARE @ID bigint

    INSERT INTO dbo.Asset (AssetType, Name)
    VALUES ('Crypto', @Name)

    SET @ID = SCOPE_IDENTITY()

    INSERT INTO dbo.CryptoAsset (ID, Ticker)
    VALUES (@ID, @Ticker)

    SELECT @ID
END

CREATE PROCEDURE dbo.GetCryptoAssetByTicker
    @Ticker nvarchar(255)
AS
BEGIN
    SET NOCOUNT ON;

    SELECT
        A.ID,
        A.Name,
        C.Ticker
    FROM dbo.Asset A
    INNER JOIN dbo.CryptoAsset C ON C.ID = A.ID
    WHERE C.Ticker = @Ticker
END

CREATE PROCEDURE dbo.ListCryptoAssets
AS
BEGIN
    SET NOCOUNT ON;

    SELECT
        A.ID,
        A.Name,
        C.Ticker
    FROM dbo.Asset A
    INNER JOIN dbo.CryptoAsset C ON C.ID = A.ID
END