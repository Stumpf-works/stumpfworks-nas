import { useState, useEffect } from 'react';
import { twofaApi, TwoFAStatus, TwoFASetupResponse } from '@/api/twofa';
import './TwoFactorAuth.css';

export function TwoFactorAuth() {
  const [status, setStatus] = useState<TwoFAStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [setupData, setSetupData] = useState<TwoFASetupResponse | null>(null);
  const [showSetup, setShowSetup] = useState(false);
  const [showDisable, setShowDisable] = useState(false);
  const [showRegenerateBackup, setShowRegenerateBackup] = useState(false);
  const [verificationCode, setVerificationCode] = useState('');
  const [newBackupCodes, setNewBackupCodes] = useState<string[]>([]);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(
    null
  );

  useEffect(() => {
    loadStatus();
  }, []);

  const loadStatus = async () => {
    try {
      const data = await twofaApi.getStatus();
      setStatus(data);
    } catch (error: any) {
      console.error('Failed to load 2FA status:', error);
    }
  };

  const handleSetup = async () => {
    try {
      setLoading(true);
      const data = await twofaApi.setup();
      setSetupData(data);
      setShowSetup(true);
      showMessage('success', '2FA setup initiated. Scan the QR code with your authenticator app.');
    } catch (error: any) {
      showMessage('error', `Failed to setup 2FA: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleEnable = async () => {
    if (!verificationCode) {
      showMessage('error', 'Please enter the verification code');
      return;
    }

    try {
      setLoading(true);
      await twofaApi.enable(verificationCode);
      showMessage('success', 'Two-factor authentication enabled successfully!');
      setShowSetup(false);
      setSetupData(null);
      setVerificationCode('');
      loadStatus();
    } catch (error: any) {
      showMessage('error', `Failed to enable 2FA: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleDisable = async () => {
    if (!verificationCode) {
      showMessage('error', 'Please enter the verification code');
      return;
    }

    try {
      setLoading(true);
      await twofaApi.disable(verificationCode);
      showMessage('success', 'Two-factor authentication disabled successfully');
      setShowDisable(false);
      setVerificationCode('');
      loadStatus();
    } catch (error: any) {
      showMessage('error', `Failed to disable 2FA: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleRegenerateBackupCodes = async () => {
    if (!verificationCode) {
      showMessage('error', 'Please enter the verification code');
      return;
    }

    try {
      setLoading(true);
      const codes = await twofaApi.regenerateBackupCodes(verificationCode);
      setNewBackupCodes(codes);
      showMessage(
        'success',
        'New backup codes generated. Please save them in a secure location!'
      );
      setVerificationCode('');
      loadStatus();
    } catch (error: any) {
      showMessage('error', `Failed to regenerate backup codes: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const showMessage = (type: 'success' | 'error', text: string) => {
    setMessage({ type, text });
    setTimeout(() => setMessage(null), 5000);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    showMessage('success', 'Copied to clipboard');
  };

  const downloadBackupCodes = (codes: string[]) => {
    const content = codes.join('\n');
    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'stumpfworks-backup-codes.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    showMessage('success', 'Backup codes downloaded');
  };

  return (
    <div className="twofa-container">
      <div className="twofa-header">
        <h3>Two-Factor Authentication</h3>
        <div className={`twofa-status ${status?.enabled ? 'enabled' : 'disabled'}`}>
          {status?.enabled ? '🔒 Enabled' : '🔓 Disabled'}
        </div>
      </div>

      {message && <div className={`message message-${message.type}`}>{message.text}</div>}

      <div className="twofa-info">
        <p>
          Two-factor authentication adds an extra layer of security to your account. When enabled,
          you'll need to provide a verification code from your authenticator app in addition to
          your password when logging in.
        </p>
      </div>

      {!status?.enabled ? (
        <div className="twofa-disabled">
          <button className="btn btn-primary" onClick={handleSetup} disabled={loading}>
            Enable Two-Factor Authentication
          </button>
        </div>
      ) : (
        <div className="twofa-enabled">
          <div className="twofa-stats">
            <div className="stat-item">
              <span className="stat-label">Status</span>
              <span className="stat-value">Active</span>
            </div>
            <div className="stat-item">
              <span className="stat-label">Backup Codes Remaining</span>
              <span className="stat-value">{status.backupCodesRemaining}</span>
            </div>
          </div>

          <div className="twofa-actions">
            <button
              className="btn"
              onClick={() => setShowRegenerateBackup(true)}
              disabled={loading}
            >
              Regenerate Backup Codes
            </button>
            <button
              className="btn btn-danger"
              onClick={() => setShowDisable(true)}
              disabled={loading}
            >
              Disable 2FA
            </button>
          </div>
        </div>
      )}

      {/* Setup Dialog */}
      {showSetup && setupData && (
        <div className="dialog-overlay" onClick={() => setShowSetup(false)}>
          <div className="dialog dialog-large" onClick={(e) => e.stopPropagation()}>
            <div className="dialog-header">
              <h3>Enable Two-Factor Authentication</h3>
              <button className="dialog-close" onClick={() => setShowSetup(false)}>
                ×
              </button>
            </div>
            <div className="dialog-body">
              <div className="setup-steps">
                <div className="setup-step">
                  <div className="step-number">1</div>
                  <div className="step-content">
                    <h4>Scan QR Code</h4>
                    <p>
                      Use your authenticator app (Google Authenticator, Authy, 1Password, etc.) to
                      scan this QR code:
                    </p>
                    <div className="qr-code-container">
                      <img
                        src={`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(
                          setupData.qrCodeUrl
                        )}`}
                        alt="QR Code"
                        className="qr-code"
                      />
                    </div>
                    <div className="manual-entry">
                      <p>Or enter this code manually:</p>
                      <div className="secret-code">
                        <code>{setupData.secret}</code>
                        <button
                          className="btn btn-sm"
                          onClick={() => copyToClipboard(setupData.secret)}
                        >
                          Copy
                        </button>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="setup-step">
                  <div className="step-number">2</div>
                  <div className="step-content">
                    <h4>Save Backup Codes</h4>
                    <p>
                      Save these backup codes in a secure location. You can use them to access your
                      account if you lose your authenticator device.
                    </p>
                    <div className="backup-codes">
                      {setupData.backupCodes.map((code, idx) => (
                        <div key={idx} className="backup-code">
                          {code}
                        </div>
                      ))}
                    </div>
                    <button
                      className="btn btn-sm"
                      onClick={() => downloadBackupCodes(setupData.backupCodes)}
                    >
                      Download Backup Codes
                    </button>
                  </div>
                </div>

                <div className="setup-step">
                  <div className="step-number">3</div>
                  <div className="step-content">
                    <h4>Verify Setup</h4>
                    <p>Enter the 6-digit code from your authenticator app to complete setup:</p>
                    <input
                      type="text"
                      className="verification-input"
                      placeholder="000000"
                      maxLength={6}
                      value={verificationCode}
                      onChange={(e) =>
                        setVerificationCode(e.target.value.replace(/\D/g, '').slice(0, 6))
                      }
                      onKeyPress={(e) => {
                        if (e.key === 'Enter' && verificationCode.length === 6) {
                          handleEnable();
                        }
                      }}
                    />
                  </div>
                </div>
              </div>
            </div>
            <div className="dialog-footer">
              <button className="btn" onClick={() => setShowSetup(false)}>
                Cancel
              </button>
              <button
                className="btn btn-primary"
                onClick={handleEnable}
                disabled={loading || verificationCode.length !== 6}
              >
                Enable 2FA
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Disable Dialog */}
      {showDisable && (
        <div className="dialog-overlay" onClick={() => setShowDisable(false)}>
          <div className="dialog" onClick={(e) => e.stopPropagation()}>
            <div className="dialog-header">
              <h3>Disable Two-Factor Authentication</h3>
              <button className="dialog-close" onClick={() => setShowDisable(false)}>
                ×
              </button>
            </div>
            <div className="dialog-body">
              <div className="warning-message">
                <p>
                  ⚠️ Disabling two-factor authentication will make your account less secure. Are
                  you sure you want to continue?
                </p>
              </div>
              <div className="form-group">
                <label>Enter your verification code to confirm:</label>
                <input
                  type="text"
                  className="verification-input"
                  placeholder="000000"
                  maxLength={6}
                  value={verificationCode}
                  onChange={(e) =>
                    setVerificationCode(e.target.value.replace(/\D/g, '').slice(0, 6))
                  }
                  onKeyPress={(e) => {
                    if (e.key === 'Enter' && verificationCode.length === 6) {
                      handleDisable();
                    }
                  }}
                />
              </div>
            </div>
            <div className="dialog-footer">
              <button className="btn" onClick={() => setShowDisable(false)}>
                Cancel
              </button>
              <button
                className="btn btn-danger"
                onClick={handleDisable}
                disabled={loading || verificationCode.length !== 6}
              >
                Disable 2FA
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Regenerate Backup Codes Dialog */}
      {showRegenerateBackup && (
        <div className="dialog-overlay" onClick={() => setShowRegenerateBackup(false)}>
          <div className="dialog" onClick={(e) => e.stopPropagation()}>
            <div className="dialog-header">
              <h3>Regenerate Backup Codes</h3>
              <button
                className="dialog-close"
                onClick={() => setShowRegenerateBackup(false)}
              >
                ×
              </button>
            </div>
            <div className="dialog-body">
              {newBackupCodes.length === 0 ? (
                <>
                  <div className="warning-message">
                    <p>
                      ⚠️ This will invalidate your old backup codes and generate new ones. Make
                      sure to save the new codes in a secure location.
                    </p>
                  </div>
                  <div className="form-group">
                    <label>Enter your verification code to confirm:</label>
                    <input
                      type="text"
                      className="verification-input"
                      placeholder="000000"
                      maxLength={6}
                      value={verificationCode}
                      onChange={(e) =>
                        setVerificationCode(e.target.value.replace(/\D/g, '').slice(0, 6))
                      }
                      onKeyPress={(e) => {
                        if (e.key === 'Enter' && verificationCode.length === 6) {
                          handleRegenerateBackupCodes();
                        }
                      }}
                    />
                  </div>
                </>
              ) : (
                <div className="backup-codes-result">
                  <p>Your new backup codes:</p>
                  <div className="backup-codes">
                    {newBackupCodes.map((code, idx) => (
                      <div key={idx} className="backup-code">
                        {code}
                      </div>
                    ))}
                  </div>
                  <button
                    className="btn btn-primary"
                    onClick={() => downloadBackupCodes(newBackupCodes)}
                  >
                    Download Backup Codes
                  </button>
                </div>
              )}
            </div>
            <div className="dialog-footer">
              {newBackupCodes.length === 0 ? (
                <>
                  <button className="btn" onClick={() => setShowRegenerateBackup(false)}>
                    Cancel
                  </button>
                  <button
                    className="btn btn-primary"
                    onClick={handleRegenerateBackupCodes}
                    disabled={loading || verificationCode.length !== 6}
                  >
                    Regenerate Codes
                  </button>
                </>
              ) : (
                <button
                  className="btn"
                  onClick={() => {
                    setShowRegenerateBackup(false);
                    setNewBackupCodes([]);
                  }}
                >
                  Close
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
