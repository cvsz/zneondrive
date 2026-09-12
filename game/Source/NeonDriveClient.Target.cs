using UnrealBuildTool;
using System.Collections.Generic;

public class NeonDriveClientTarget : TargetRules
{
    public NeonDriveClientTarget(TargetInfo Target) : base(Target)
    {
        Type = TargetType.Client;
        DefaultBuildSettings = BuildSettingsVersion.V5;
        IncludeOrderVersion = EngineIncludeOrderVersion.Latest;
        ExtraModuleNames.Add("NeonDrive");
    }
}
