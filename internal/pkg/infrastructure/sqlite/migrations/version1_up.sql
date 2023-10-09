--- Release 1
----- Tables
CREATE TABLE IF NOT EXISTS Asset
(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    AssetType TEXT NOT NULL CHECK (AssetType IN ('Stock', 'Crypto')),
    Name TEXT NOT NULL,
    IsActive BIT NOT NULL DEFAULT (1)
);

CREATE TABLE IF NOT EXISTS StockAsset
(
    ID INTEGER PRIMARY KEY NOT NULL,
    Ticker TEXT NOT NULL,
    Exchange TEXT NOT NULL,
    Currency TEXT NOT NULL,
    FOREIGN KEY (ID) REFERENCES Asset (ID)
);

CREATE TABLE IF NOT EXISTS CryptoAsset
(
    ID INTEGER PRIMARY KEY NOT NULL,
    Ticker TEXT NOT NULL,
    FOREIGN KEY (ID) REFERENCES Asset (ID)
);

CREATE TABLE IF NOT EXISTS AssetPrice
(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    AssetID INTEGER NOT NULL,
    Price DECIMAL(10, 2) NOT NULL,
    Currency TEXT NOT NULL,
    Date DATE NOT NULL,
    FOREIGN KEY (AssetID) REFERENCES Asset (ID)
);

CREATE TABLE IF NOT EXISTS Ledger (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id VARCHAR(255) NOT NULL,
    exchange VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    asset_type VARCHAR(255) NOT NULL CHECK (asset_type IN ('CRYPTO', 'STOCK')),
    asset_id INTEGER NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    cost_basis DECIMAL(10, 2) NOT NULL,
    cost_basis_currency VARCHAR(3) NOT NULL,
    fees DECIMAL(10, 2) NOT NULL,
    fees_currency VARCHAR(3) NOT NULL,
    spent BIT NOT NULL DEFAULT 0,
    spent_at DATE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (asset_id) REFERENCES Asset (id)
);

----- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS IX_StockAsset_Ticker ON StockAsset (Ticker);
CREATE UNIQUE INDEX IF NOT EXISTS IX_CryptoAsset_Ticker ON CryptoAsset (Ticker);
CREATE UNIQUE INDEX IF NOT EXISTS IX_AssetPrice_AssetID_Date ON AssetPrice (AssetID, Date);
CREATE UNIQUE INDEX IF NOT EXISTS IX_Ledger_TransactionID ON Ledger (transaction_id);
CREATE INDEX IF NOT EXISTS IX_Ledger_AssetID ON Ledger (asset_id);
CREATE INDEX IF NOT EXISTS IX_Ledger_Spent ON Ledger (spent);
CREATE INDEX IF NOT EXISTS IX_Ledger_AssetType ON Ledger (asset_type);
