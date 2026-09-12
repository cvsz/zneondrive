using UnrealBuildTool;
using System.Collections.Generic;

public class NeonDriveServerTarget : TargetRules
{
    public NeonDriveServerTarget(TargetInfo Target) : base(Target)
    {
        Type = TargetType.Server;
        DefaultBuildSettings = BuildSettingsVersion.V5;
        IncludeOrderVersion = EngineIncludeOrderVersion.Latest;
        ExtraModuleNames.Add("NeonDrive");
    }
}
