using UnrealBuildTool;
using System.Collections.Generic;

public class NeonDriveTarget : TargetRules
{
    public NeonDriveTarget(TargetInfo Target) : base(Target)
    {
        Type = TargetType.Game;
        DefaultBuildSettings = BuildSettingsVersion.V7;
        IncludeOrderVersion = EngineIncludeOrderVersion.Latest;
        ExtraModuleNames.Add("NeonDrive");
    }
}
